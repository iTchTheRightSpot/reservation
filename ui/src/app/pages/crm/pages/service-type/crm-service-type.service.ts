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
import { ServicePayToStaffload } from '@crm/pages/staff/pages/all/ui/shared/link-service/link-service.model';
import { Cache } from '@shared/data-access/cache';

@Injectable({
  providedIn: 'root'
})
export class CRMServiceTypeService {
  private static readonly serviceTypesByStaffCache = new Cache<
    string,
    string[]
  >();

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  private readonly allServiceTypesCache = new BehaviorSubject<
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
      : this.allServiceTypesCache.asObservable().pipe(
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
                      this.allServiceTypesCache.next(arr);
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
      map(() => {
        this.allServiceTypesCache.next(undefined);
        return { state: ApiState.LOADED } as ApiResponse<T>;
      }),
      startWith({ state: ApiState.LOADING } as ApiResponse<T>),
      catchError(e => {
        this.toast.message({ state: ToastEnum.ERROR, message: e.message });
        return of(err<T>(e));
      })
    );
  };

  readonly linkServiceToStaff = (
    o: ServicePayToStaffload
  ): Observable<ApiResponse<any>> =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<any>>({ state: ApiState.LOADING }),
              timer(1000).pipe(
                concatMap(() =>
                  of<ApiResponse<any>>({ state: ApiState.LOADED })
                )
              )
            )
          )
        )
      : this.http
          .post<ServicePayToStaffload>(
            `${environment.domain}crm/service/staff`,
            o,
            {
              withCredentials: true
            }
          )
          .pipe(
            map(() => {
              this.toast.message({
                message: 'linked service to staff!',
                state: ToastEnum.SUCCESS
              });
              return { state: ApiState.LOADED } as ApiResponse<any>;
            }),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => {
              this.toast.message({
                state: ToastEnum.ERROR,
                message: e.message
              });
              return of(err<any>(e));
            })
          );

  readonly servicesByStaff = (
    staffId: string
  ): Observable<ApiResponse<string[]>> => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<string[]>>({
              state: ApiState.LOADING
            }),
            timer(2000).pipe(
              concatMap(() =>
                of<ApiResponse<string[]>>({
                  state: ApiState.LOADED,
                  data: ['lawn', 'pedicure', 'mens hair']
                })
              )
            )
          )
        )
      );

    const req$ = this.http
      .get<
        string[]
      >(`${environment.domain}crm/services/staff`, { withCredentials: true })
      .pipe(
        map(arr => {
          CRMServiceTypeService.serviceTypesByStaffCache.setItem(staffId, arr);
          return {
            state: ApiState.LOADED,
            data: arr
          } as ApiResponse<string[]>;
        }),
        startWith({ state: ApiState.LOADING } as ApiResponse<string[]>),
        catchError(e => of(err<string[]>(e)))
      );

    return CRMServiceTypeService.serviceTypesByStaffCache
      .getItem(staffId)
      .pipe(
        switchMap(arr =>
          arr
            ? of<ApiResponse<string[]>>({ state: ApiState.LOADED, data: arr })
            : arr === undefined
              ? req$
              : of<ApiResponse<string[]>>({ state: ApiState.LOADED, data: [] })
        )
      );
  };

  readonly deLinkServiceFromStaff = (o: ServicePayToStaffload) =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<any>>({ state: ApiState.LOADING }),
              timer(1000).pipe(
                concatMap(() =>
                  of<ApiResponse<any>>({ state: ApiState.LOADED })
                )
              )
            )
          )
        )
      : this.http
          .delete<
            ApiResponse<any>
          >(`${environment.domain}crm/service/staff?staff_id=${o.staff_id}&service=${o.service}`, { withCredentials: true })
          .pipe(
            map(() => {
              this.toast.message({
                message: 'de-linked service from staff',
                state: ToastEnum.SUCCESS
              });
              return { state: ApiState.LOADED } as ApiResponse<any>;
            }),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => {
              this.toast.message({
                state: ToastEnum.ERROR,
                message: e.message
              });
              return of(err<any>(e));
            })
          );
}
