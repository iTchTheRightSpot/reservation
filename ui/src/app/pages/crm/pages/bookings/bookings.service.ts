import { inject, Injectable } from '@angular/core';
import { Cache } from '@shared/data-access/cache';
import {
  BookingsModel,
  BookingsRequestPayload,
  DummyBookingsModels,
  UpdateBookingStatusPayload
} from './bookings.model';
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
import { HttpClient, HttpParams } from '@angular/common/http';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { err, TIMEZONE } from '@root/app.util';
import {
  ConfirmModel,
  DateModel,
  DummyDates,
  DummyStaffModels,
  StaffModel
} from '@shared/model/shared.model';

@Injectable({
  providedIn: 'root'
})
export class BookingsService {
  private static readonly BookingsCache = new Cache<string, BookingsModel[]>();
  private static readonly StaffsByServicesCache = new Cache<
    string,
    StaffModel[]
  >();
  private static readonly ValidDatesCache = new Cache<string, DateModel[]>();

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly bookings = (o: BookingsRequestPayload) => {
    if (!environment.production)
      return of<ApiResponse<BookingsModel[]>>({
        state: ApiState.LOADED,
        data: DummyBookingsModels(50)
      });

    const k = `${o.user_id}_${o.date.getMonth()}_${o.date.getFullYear()}`;
    return BookingsService.BookingsCache.getItem(k).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<BookingsModel[]>>({
              state: ApiState.LOADED,
              data: arr
            })
          : arr === null
            ? of<ApiResponse<BookingsModel[]>>({
                state: ApiState.LOADED,
                data: []
              })
            : this.http
                .get<
                  BookingsModel[]
                >(`${environment.domain}crm/reservation?user_id=${o.user_id}&month=${1 + o.date.getMonth()}&year=${o.date.getFullYear()}&timezone=${TIMEZONE}`, { withCredentials: true })
                .pipe(
                  map(arr => {
                    BookingsService.BookingsCache.setItem(k, arr);
                    return { state: ApiState.LOADED, data: arr };
                  }),
                  startWith({ state: ApiState.LOADING } as ApiResponse<
                    BookingsModel[]
                  >),
                  catchError(e => {
                    this.toast.message({
                      state: ToastEnum.ERROR,
                      message: e.message
                    });
                    return of(err<BookingsModel[]>(e));
                  })
                )
      )
    );
  };

  readonly updateBookingStatus = (u: UpdateBookingStatusPayload) =>
    !environment.production
      ? of<ApiResponse<any>>({ state: ApiState.LOADED })
      : this.http
          .put<any>(`${environment.domain}crm/reservation`, u, {
            withCredentials: true
          })
          .pipe(
            map(() => {
              BookingsService.BookingsCache.clear();
              return { state: ApiState.LOADED };
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

  readonly staffsByServiceTypes = (services: string[]) => {
    if (!environment.production)
      return of<ApiResponse<StaffModel[]>>({
        state: ApiState.LOADED,
        data: DummyStaffModels
      });

    const key = services.join('_');

    return BookingsService.StaffsByServicesCache.getItem(key).pipe(
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
              BookingsService.StaffsByServicesCache.setItem(key, staffs)
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

    return BookingsService.ValidDatesCache.getItem(key).pipe(
      switchMap(arr => {
        if (arr)
          return of<ApiResponse<DateModel[]>>({
            state: ApiState.LOADED,
            data: arr
          });

        let params = new HttpParams();
        params = params.append('staff_id', o.staff_id);
        params = params.append('day', o.date.getDate());
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
            tap(arr => BookingsService.ValidDatesCache.setItem(key, arr)),
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
              BookingsService.BookingsCache.clear();
              BookingsService.ValidDatesCache.clear();
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
