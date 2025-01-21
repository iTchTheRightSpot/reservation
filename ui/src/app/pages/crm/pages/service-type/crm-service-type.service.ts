import { inject, Injectable } from '@angular/core';
import {
  CRM_DummyServiceTypes,
  CRM_ServiceTypeModel
} from './crm-service-type.model';
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
  timer
} from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';

@Injectable({
  providedIn: 'root'
})
export class CRMServiceTypeService {
  private readonly http = inject(HttpClient);

  private readonly cache = new BehaviorSubject<
    CRM_ServiceTypeModel[] | undefined
  >(undefined);

  readonly all = () =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<CRM_ServiceTypeModel[]>>({
                state: ApiState.LOADING
              }),
              timer(2000).pipe(
                concatMap(() =>
                  of<ApiResponse<CRM_ServiceTypeModel[]>>({
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
              ? of<ApiResponse<CRM_ServiceTypeModel[]>>({
                  state: ApiState.LOADED,
                  data: arr
                })
              : this.http
                  .get<
                    CRM_ServiceTypeModel[]
                  >(`${environment.domain}crm/service`, { withCredentials: true })
                  .pipe(
                    map(arr => {
                      this.cache.next(arr);
                      return {
                        state: ApiState.LOADED,
                        data: arr
                      } as ApiResponse<CRM_ServiceTypeModel[]>;
                    }),
                    startWith({ state: ApiState.LOADING } as ApiResponse<
                      CRM_ServiceTypeModel[]
                    >),
                    catchError(e => of(err<CRM_ServiceTypeModel[]>(e)))
                  )
          )
        );

  readonly create = (m: CRM_ServiceTypeModel) =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<CRM_ServiceTypeModel[]>>({
                state: ApiState.LOADING
              }),
              timer(2000).pipe(
                concatMap(() =>
                  of<ApiResponse<CRM_ServiceTypeModel[]>>({
                    state: ApiState.LOADED,
                    data: CRM_DummyServiceTypes(50)
                  })
                )
              )
            )
          )
        )
      : this.http
          .post<any>(`${environment.domain}crm/service`, m, {
            withCredentials: true
          })
          .pipe(
            map(() => {
              this.cache.next(undefined);
              return { state: ApiState.LOADED } as ApiResponse<any>;
            }),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => of(err<any>(e)))
          );
}
