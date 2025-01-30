import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '@env/environment';
import {
  BehaviorSubject,
  catchError,
  concat,
  concatMap,
  delay,
  map,
  merge,
  of,
  startWith,
  switchMap,
  tap,
  timer
} from 'rxjs';
import { ActiveUser, ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { toSignal } from '@angular/core/rxjs-interop';
import { Permission, Role } from '@crm/pages/staff/pages/all/crm-staff.model';
import { ToastEnum, ToastService } from './toast.service';

export interface LoginModel {
  email: string;
  password: string;
}

export interface RegisterModel {
  firstname: string;
  lastname: string;
  email: string;
  password: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  private readonly cache = new BehaviorSubject<ActiveUser | undefined>(
    undefined
  );

  private readonly activeUserRequest = () =>
    !environment.production
      ? of<ActiveUser>({
          user_id: '1',
          firstname: 'test user',
          image_key: null,
          access_controls: [
            {
              role: Role.STAFF,
              permissions: [Permission.READ, Permission.WRITE]
            }
          ]
        })
      : this.http
          .get<ActiveUser>(`${environment.domain}active`, {
            withCredentials: true
          })
          .pipe(
            tap(o => {
              if (o && Object.keys(o).length > 0) this.cache.next(o);
            }),
            catchError(e => {
              const m = e.error ? e.error.message : e.message;
              console.error(m);
              return of(m);
            })
          );

  readonly activeUser = toSignal(
    this.cache
      .asObservable()
      .pipe(switchMap(o => (o ? of<ActiveUser>(o) : this.activeUserRequest()))),
    { initialValue: undefined }
  );

  readonly login = (obj: LoginModel) =>
    environment.production
      ? this.http
          .post<any>(`${environment.domain}account/login`, obj, {
            withCredentials: true
          })
          .pipe(
            switchMap(() =>
              this.activeUserRequest().pipe(
                map(() => ({ state: ApiState.LOADED }) as ApiResponse<any>)
              )
            ),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => of(err<any>(e)))
          )
      : of(obj).pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<any>>({ state: ApiState.LOADING }),
              timer(2000).pipe(
                concatMap(() =>
                  of<ApiResponse<any>>({ state: ApiState.LOADED })
                )
              )
            )
          )
        );

  readonly register = (obj: RegisterModel) =>
    environment.production
      ? this.http
          .post<any>(`${environment.domain}account/register`, obj, {
            withCredentials: true
          })
          .pipe(
            map(() => ({ state: ApiState.LOADED }) as ApiResponse<any>),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => of(err<any>(e)))
          )
      : of('register').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<any>>({ state: ApiState.LOADING }),
              timer(2000).pipe(
                concatMap(() =>
                  of<ApiResponse<any>>({ state: ApiState.LOADED })
                )
              )
            )
          )
        );

  readonly logout = () =>
    environment.production
      ? this.http
          .post<
            ApiResponse<any>
          >(`${environment.domain}logout`, {}, { withCredentials: true })
          .pipe(
            map(() => <ApiResponse<any>>{ state: ApiState.LOADED }),
            startWith(<ApiResponse<any>>{ state: ApiState.LOADING }),
            catchError(e => {
              this.toast.message({
                message: e.message,
                state: ToastEnum.ERROR
              });
              return of(err<any>(e));
            })
          )
      : merge(
          of({ state: ApiState.LOADING }),
          of({ state: ApiState.LOADED }).pipe(delay(2000))
        );
}
