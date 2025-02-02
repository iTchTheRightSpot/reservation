import { inject, Injectable } from '@angular/core';
import { Cache } from '@shared/data-access/cache';
import {
  BookingsModel,
  BookingsRequestPayload,
  DummyBookingsModels,
  UpdateBookingStatusPayload
} from './bookings.model';
import { environment } from '@env/environment';
import { catchError, map, of, startWith, switchMap } from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { HttpClient } from '@angular/common/http';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { err, TIMEZONE } from '@root/app.util';

@Injectable({
  providedIn: 'root'
})
export class BookingsService {
  private static readonly cache = new Cache<string, BookingsModel[]>();

  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly bookings = (o: BookingsRequestPayload) => {
    if (!environment.production)
      return of<ApiResponse<BookingsModel[]>>({
        state: ApiState.LOADED,
        data: DummyBookingsModels(50)
      });

    const k = `${o.user_id}_${o.date.getMonth()}_${o.date.getFullYear()}`;
    return BookingsService.cache.getItem(k).pipe(
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
                    BookingsService.cache.setItem(k, arr);
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
              BookingsService.cache.clear();
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
}
