import { inject, Injectable } from '@angular/core';
import {
  CRM_DummyServiceTypes,
  CRMServiceTypeModel
} from './crm-service-type.model';
import { HttpClient } from '@angular/common/http';
import { environment } from '@env/environment';
import {
  BehaviorSubject,
  catchError,
  concat,
  concatMap,
  map,
  Observable,
  of,
  startWith,
  switchMap,
  timer
} from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';

@Injectable({
  providedIn: 'root'
})
export class CRMServiceTypeService {
  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  private readonly cache = new BehaviorSubject<
    CRMServiceTypeModel[] | undefined
  >(undefined);

  readonly all = (): Observable<ApiResponse<CRMServiceTypeModel[]>> =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<CRMServiceTypeModel[]>>({
                state: ApiState.LOADING
              }),
              timer(2000).pipe(
                concatMap(() =>
                  of<ApiResponse<CRMServiceTypeModel[]>>({
                    state: ApiState.LOADED,
                    data: CRM_DummyServiceTypes(50)
                  })
                )
              )
            )
          )
        )
      : this.cache.asObservable().pipe(
          switchMap(arr =>
            arr
              ? of<ApiResponse<CRMServiceTypeModel[]>>({
                  state: ApiState.LOADED,
                  data: arr
                })
              : this.http
                  .get<
                    CRMServiceTypeModel[]
                  >(`${environment.domain}crm/services`, { withCredentials: true })
                  .pipe(
                    map(arr => {
                      this.cache.next(arr);
                      return {
                        state: ApiState.LOADED,
                        data: arr
                      } as ApiResponse<CRMServiceTypeModel[]>;
                    }),
                    startWith({ state: ApiState.LOADING } as ApiResponse<
                      CRMServiceTypeModel[]
                    >),
                    catchError(e => of(err<CRMServiceTypeModel[]>(e)))
                  )
          )
        );

  readonly create = (m: CRMServiceTypeModel) => this.write<any>('POST', m);

  readonly update = (m: CRMServiceTypeModel) => this.write<any>('PUT', m);

  private readonly write = <T>(
    method: 'POST' | 'PUT',
    body: CRMServiceTypeModel
  ) => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<T>>({ state: ApiState.LOADING }),
            timer(2000).pipe(
              concatMap(() => of<ApiResponse<T>>({ state: ApiState.LOADED }))
            )
          )
        )
      );

    const url = `${environment.domain}crm/service`;

    const req$ =
      method === 'POST'
        ? this.http.post<T>(url, body, { withCredentials: true })
        : this.http.put<T>(url, body, { withCredentials: true });

    return req$.pipe(
      map(res => {
        this.cache.next(undefined);
        return { state: ApiState.LOADED } as ApiResponse<T>;
      }),
      startWith({ state: ApiState.LOADING } as ApiResponse<T>),
      catchError(e => {
        this.toast.message({ state: ToastEnum.ERROR, message: e.message });
        return of(err<T>(e));
      })
    );
  };
}
