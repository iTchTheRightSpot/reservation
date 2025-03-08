import { inject, Injectable } from '@angular/core';
import { environment } from '@env/environment';
import {
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
import { ApiResponse, ApiState } from '@root/app.model';
import {
  ConfirmModel,
  DateModel,
  DummyDates,
  DummyStaffModels,
  StaffModel
} from '@shared/model/shared.model';
import { HttpClient, HttpParams } from '@angular/common/http';
import { ToastEnum, ToastService } from './toast.service';
import { err, TIMEZONE } from '@root/app.util';
import { Cache } from './cache';

@Injectable({
  providedIn: 'root'
})
export class BookingService {
  private static readonly StaffsByServicesCache = new Cache<
    string,
    StaffModel[]
  >();
  private static readonly ValidDatesCache = new Cache<string, DateModel[]>();

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly staffsByServiceTypes = (services: string[]) => {
    if (!environment.production)
      return of<ApiResponse<StaffModel[]>>({
        state: ApiState.LOADED,
        data: DummyStaffModels
      });

    const key = services.join('_');

    return BookingService.StaffsByServicesCache.getItem(key).pipe(
      switchMap(value => {
        if (value)
          return of<ApiResponse<StaffModel[]>>({
            state: ApiState.LOADED,
            data: value
          });

        let params = new HttpParams();
        services.forEach(
          service => (params = params.append('name', service.trim()))
        );

        return this.http
          .get<
            StaffModel[]
          >(`${environment.domain}service/staffs`, { withCredentials: true, params: params })
          .pipe(
            tap(staffs =>
              BookingService.StaffsByServicesCache.setItem(key, staffs)
            ),
            map(
              arr =>
                ({ state: ApiState.LOADED, data: arr }) as ApiResponse<
                  StaffModel[]
                >
            ),
            catchError(e => {
              this.toast.message({
                state: ToastEnum.ERROR,
                message: e.message
              });
              return of(err<StaffModel[]>(e));
            })
          );
      })
    );
  };

  readonly validDates = (o: {
    services: string[];
    staff_id: string;
    date: Date;
  }) => {
    if (!environment.production)
      return of<ApiResponse<DateModel[]>>({
        state: ApiState.LOADED,
        data: DummyDates(o.date.getMonth(), o.date.getFullYear())
      });

    const key = `${o.staff_id}_${o.services.join('_')}_${o.date.getMonth()}_${o.date.getFullYear()}_${TIMEZONE}`;

    return BookingService.ValidDatesCache.getItem(key).pipe(
      switchMap(arr => {
        if (arr)
          return of<ApiResponse<DateModel[]>>({
            state: ApiState.LOADED,
            data: arr
          });

        let params = new HttpParams();
        params = params.append('staff_id', o.staff_id);
        params = params.append('day', 1);
        params = params.append('month', 1 + o.date.getMonth());
        params = params.append('year', o.date.getFullYear());
        o.services.forEach(
          service => (params = params.append('service', service.trim()))
        );
        params = params.append('timezone', TIMEZONE);

        return this.http
          .get<
            DateModel[]
          >(`${environment.domain}reservation`, { withCredentials: true, params: params })
          .pipe(
            tap(arr => BookingService.ValidDatesCache.setItem(key, arr)),
            map(
              arr =>
                ({ state: ApiState.LOADED, data: arr }) as ApiResponse<
                  DateModel[]
                >
            ),
            startWith({ state: ApiState.LOADING } as ApiResponse<DateModel[]>),
            catchError(e => {
              this.toast.message({
                state: ToastEnum.ERROR,
                message: e.message
              });
              return of(err<DateModel[]>(e));
            })
          );
      })
    );
  };

  readonly create = (o: ConfirmModel) =>
    environment.production
      ? this.http
          .post<any>(`${environment.domain}reservation`, o, {
            withCredentials: true
          })
          .pipe(
            map(() => {
              BookingService.ValidDatesCache.clear();
              this.toast.message({
                state: ToastEnum.SUCCESS,
                message: 'booking created!'
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
