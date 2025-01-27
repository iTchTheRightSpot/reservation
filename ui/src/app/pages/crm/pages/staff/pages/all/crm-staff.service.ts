import { inject, Injectable } from '@angular/core';
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
import { CRMStaffModel, DummyCRMStaffModels } from './crm-staff.model';
import { HttpClient } from '@angular/common/http';
import { environment } from '@env/environment';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';

@Injectable({
  providedIn: 'root'
})
export class CRMStaffsService {
  private readonly cache = new BehaviorSubject<
    CRMStaffModel[] | null | undefined
  >(null);

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly staffs = () =>
    !environment.production
      ? of('yes').pipe(
          concatMap(() =>
            concat(
              of<ApiResponse<CRMStaffModel[]>>({ state: ApiState.LOADING }),

              timer(1000).pipe(
                concatMap(() =>
                  of<ApiResponse<CRMStaffModel[]>>({
                    state: ApiState.LOADED,
                    data: DummyCRMStaffModels(25)
                  })
                )
              )
            )
          )
        )
      : this.cache.asObservable().pipe(
          switchMap(arr =>
            arr === null
              ? this.allStaffsRequest()
              : arr === undefined
                ? of<ApiResponse<CRMStaffModel[]>>({
                    state: ApiState.LOADED,
                    data: []
                  })
                : of<ApiResponse<CRMStaffModel[]>>({
                    state: ApiState.LOADED,
                    data: arr
                  })
          )
        );

  private readonly allStaffsRequest = () =>
    this.http
      .get<
        CRMStaffModel[]
      >(`${environment.domain}staffs`, { withCredentials: true })
      .pipe(
        map(arr => {
          this.cache.next(arr);
          return { state: ApiState.LOADED, data: arr } as ApiResponse<
            CRMStaffModel[]
          >;
        }),
        startWith({ state: ApiState.LOADING } as ApiResponse<CRMStaffModel[]>),
        catchError(e => {
          this.toast.message({ state: ToastEnum.ERROR, message: e.message });
          return of(err<CRMStaffModel[]>(e));
        })
      );
}
