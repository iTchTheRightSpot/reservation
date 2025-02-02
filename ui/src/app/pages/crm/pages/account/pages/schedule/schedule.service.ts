import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Cache } from '@shared/data-access/cache';
import {
  DummySchedules,
  Schedule,
  CreateScheduleModel,
  UpdateScheduleModel
} from './schedule.model';
import { environment } from '@env/environment';
import {
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
import { err, TIMEZONE } from '@root/app.util';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { StaffScheduleEmitter } from '@crm/pages/staff/pages/all/ui/shared/staff-schedule/staff-schedule.model';
import moment from 'moment-timezone';

@Injectable({
  providedIn: 'root'
})
export class ScheduleService {
  static readonly AllSchedulesCache = new Cache<string, Schedule[]>();
  static readonly SchedulesByStaffCache = new Cache<string, Schedule[]>();

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly all = (d: Date, page: number, size: number) => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<Schedule[]>>({
              state: ApiState.LOADING
            }),
            timer(0).pipe(
              concatMap(() =>
                of<ApiResponse<Schedule[]>>({
                  state: ApiState.LOADED,
                  data: DummySchedules(50)
                })
              )
            )
          )
        )
      );

    const key = `${d.getMonth()}_${d.getFullYear()}_${page}_${size}`;
    return ScheduleService.AllSchedulesCache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<Schedule[]>>({ state: ApiState.LOADED, data: arr })
          : this.http
              .get<
                Schedule[]
              >(`${environment.domain}schedules?month=${1 + d.getMonth()}&year=${d.getFullYear()}&page=${page}&size=${size}`, { withCredentials: true })
              .pipe(
                map(arr => {
                  ScheduleService.AllSchedulesCache.setItem(key, arr);
                  return {
                    state: ApiState.LOADED,
                    data: arr
                  } as ApiResponse<Schedule[]>;
                }),
                startWith({ state: ApiState.LOADING } as ApiResponse<
                  Schedule[]
                >),
                catchError(e => {
                  this.toast.message({
                    state: ToastEnum.ERROR,
                    message: e.message
                  });
                  return of(err<Schedule[]>(e));
                })
              )
      )
    );
  };

  readonly schedulesByStaff = (
    o: StaffScheduleEmitter
  ): Observable<ApiResponse<Schedule[]>> => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<Schedule[]>>({
              state: ApiState.LOADING
            }),
            timer(0).pipe(
              concatMap(() =>
                of<ApiResponse<Schedule[]>>({
                  state: ApiState.LOADED,
                  data: DummySchedules(20)
                })
              )
            )
          )
        )
      );

    const key = `${o.staff_id}_${o.date.getMonth()}_${o.date.getFullYear()}`;
    return ScheduleService.SchedulesByStaffCache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<Schedule[]>>({ state: ApiState.LOADED, data: arr })
          : this.http
              .get<
                Schedule[]
              >(`${environment.domain}schedules/staff?staff_id=${o.staff_id}&month=${1 + o.date.getMonth()}&year=${o.date.getFullYear()}`, { withCredentials: true })
              .pipe(
                map(arr => {
                  ScheduleService.SchedulesByStaffCache.setItem(key, arr);
                  return {
                    state: ApiState.LOADED,
                    data: arr
                  } as ApiResponse<Schedule[]>;
                }),
                startWith({ state: ApiState.LOADING } as ApiResponse<
                  Schedule[]
                >),
                catchError(e => {
                  this.toast.message({
                    state: ToastEnum.ERROR,
                    message: e.message
                  });
                  return of(err<Schedule[]>(e));
                })
              )
      )
    );
  };

  readonly create = (o: CreateScheduleModel) => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<any>>({ state: ApiState.LOADING }),
            timer(0).pipe(
              concatMap(() => of<ApiResponse<any>>({ state: ApiState.LOADED }))
            )
          )
        )
      );

    const obj = {
      staff_id: o.staff_id,
      schedule_segments: [
        {
          is_visible: o.is_visible,
          is_reoccurring: o.is_reoccurring,
          start: moment.tz(o.start, TIMEZONE).utc().format(),
          duration: o.duration
        }
      ]
    };

    return this.http
      .post<any>(`${environment.domain}schedule`, obj, {
        withCredentials: true
      })
      .pipe(
        map(() => {
          this.toast.message({
            message: 'schedule created',
            state: ToastEnum.SUCCESS
          });
          return { state: ApiState.LOADED };
        }),
        startWith({ state: ApiState.LOADING } as ApiResponse<any>),
        catchError(e => {
          this.toast.message({ message: e.message, state: ToastEnum.ERROR });
          return of(err<any>(e));
        })
      );
  };

  readonly update = (o: UpdateScheduleModel) => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<any>>({ state: ApiState.LOADING }),
            timer(0).pipe(
              concatMap(() => of<ApiResponse<any>>({ state: ApiState.LOADED }))
            )
          )
        )
      );

    const obj = {
      schedule_id: o.schedule_id,
      is_visible: o.is_visible,
      is_reoccurring: o.is_reoccurring
    };

    return this.http
      .put<any>(`${environment.domain}schedule`, obj, {
        withCredentials: true
      })
      .pipe(
        map(() => {
          this.toast.message({
            message: 'schedule updated',
            state: ToastEnum.SUCCESS
          });
          return { state: ApiState.LOADED };
        }),
        startWith({ state: ApiState.LOADING } as ApiResponse<any>),
        catchError(e => {
          this.toast.message({ message: e.message, state: ToastEnum.ERROR });
          return of(err<any>(e));
        })
      );
  };

  readonly delete = (scheduleId: number) => {
    if (!environment.production)
      return of('yes').pipe(
        concatMap(() =>
          concat(
            of<ApiResponse<any>>({ state: ApiState.LOADING }),
            timer(0).pipe(
              concatMap(() => of<ApiResponse<any>>({ state: ApiState.LOADED }))
            )
          )
        )
      );

    return this.http
      .delete<any>(`${environment.domain}schedule/${scheduleId}`, {
        withCredentials: true
      })
      .pipe(
        map(() => {
          this.toast.message({ message: 'deleted', state: ToastEnum.SUCCESS });
          return { state: ApiState.LOADED };
        }),
        startWith({ state: ApiState.LOADING } as ApiResponse<any>),
        catchError(e => {
          this.toast.message({ message: e.message, state: ToastEnum.ERROR });
          return of(err<any>(e));
        })
      );
  };
}
