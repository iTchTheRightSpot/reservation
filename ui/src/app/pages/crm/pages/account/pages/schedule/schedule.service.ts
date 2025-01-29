import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Cache } from '@shared/data-access/cache';
import { DummySchedules, Schedule } from './schedule.model';
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
import { err } from '@root/app.util';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { StaffScheduleEmitter } from '@crm/pages/staff/pages/all/ui/shared/staff-schedule/staff-schedule.model';
import { CreateUpdateScheduleModel } from '@crm/pages/staff/pages/all/ui/shared/staff-schedule/ui/shared/shared-create-update-schedule.model';

@Injectable({
  providedIn: 'root'
})
export class ScheduleService {
  private static readonly allSchedulesCache = new Cache<string, Schedule[]>();
  private static readonly schedulesByStaffCache = new Cache<
    string,
    Schedule[]
  >();

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
    return ScheduleService.allSchedulesCache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<Schedule[]>>({ state: ApiState.LOADED, data: arr })
          : this.http
              .get<
                Schedule[]
              >(`${environment.domain}schedules?month=${1 + d.getMonth()}&year=${d.getFullYear()}&page=${page}&size=${size}`, { withCredentials: true })
              .pipe(
                map(arr => {
                  ScheduleService.allSchedulesCache.setItem(key, arr);
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
    return ScheduleService.schedulesByStaffCache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<Schedule[]>>({ state: ApiState.LOADED, data: arr })
          : this.http
              .get<
                Schedule[]
              >(`${environment.domain}schedules/staff?staff_id=${o.staff_id}&month=${1 + o.date.getMonth()}&year=${o.date.getFullYear()}`, { withCredentials: true })
              .pipe(
                map(arr => {
                  ScheduleService.schedulesByStaffCache.setItem(key, arr);
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

  readonly create = (o: CreateUpdateScheduleModel) => of<ApiResponse<any>>();

  readonly update = (o: CreateUpdateScheduleModel) => of<ApiResponse<any>>();

  readonly delete = (scheduleId: number) => of<ApiResponse<any>>();
}
