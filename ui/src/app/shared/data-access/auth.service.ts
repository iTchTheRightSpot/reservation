import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '@env/environment';
import {
  BehaviorSubject,
  catchError,
  concat,
  concatMap,
  map,
  of,
  startWith,
  switchMap,
  tap,
  timer
} from 'rxjs';
import { ActiveUser, ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { toSignal } from '@angular/core/rxjs-interop';

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

  private readonly cache = new BehaviorSubject<ActiveUser | undefined>(
    undefined
  );

  private readonly activeUserRequest = () =>
    !environment.production
      ? of<ActiveUser>({
          user_id: '1',
          firstname: 'test user',
          image_key: null,
          roles_permission: []
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
}
