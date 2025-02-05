import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { ReservationService } from '@store/pages/reservation/reservation.service';
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
import { DateModel, DummyDates } from '@shared/data-access/shared.model';
import { ApiResponse, ApiState } from '@root/app.model';
import { Cache } from '@shared/data-access/cache';
import { err, TIMEZONE } from '@root/app.util';

@Injectable({
  providedIn: 'root'
})
export class DatesService {
  private static readonly cache = new Cache<string, DateModel[]>();

  private readonly http = inject(HttpClient);
  private readonly service = inject(ReservationService);

  readonly clearCache = () => DatesService.cache.clear();

  readonly dates = (date: Date) => {
    if (!environment.production) {
      const loading = of<ApiResponse<DateModel[]>>({ state: ApiState.LOADING });
      const loaded = timer(1000).pipe(
        switchMap(() =>
          of<ApiResponse<DateModel[]>>({
            state: ApiState.LOADED,
            data: DummyDates(date.getMonth(), date.getFullYear())
          })
        )
      );

      return of('yes').pipe(concatMap(() => concat(loading, loaded)));
    }

    const obj = this.service.reservationState();
    if (!obj.services || obj.services.length < 1 || !obj.staff)
      return of<ApiResponse<DateModel[]>>({
        state: ApiState.ERROR,
        message: 'missing: service(s) & staff'
      });

    const d = new Date(date);
    if (d.getMonth() !== new Date().getMonth()) d.setDate(1);

    const key = `${obj.staff.staff_id}_${obj.services.map(s => s.name).join('_')}_${d.getMonth()}_${d.getFullYear()}_${TIMEZONE}`;

    return DatesService.cache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<DateModel[]>>({
              state: ApiState.LOADED,
              data: arr
            })
          : this.req(
              key,
              obj!.services!.map(o => o.name),
              obj.staff!.staff_id,
              d
            )
      )
    );
  };

  private readonly req = (
    key: string,
    services: string[],
    staffId: string,
    date: Date
  ) => {
    let params = new HttpParams();
    params = params.append('staff_id', staffId);
    params = params.append('day', date.getDate());
    params = params.append('month', 1 + date.getMonth());
    params = params.append('year', date.getFullYear());
    services.forEach(
      service => (params = params.append('service', service.trim()))
    );
    params = params.append('timezone', TIMEZONE);

    return this.http
      .get<
        DateModel[]
      >(`${environment.domain}reservation`, { withCredentials: true, params: params })
      .pipe(
        tap(arr => DatesService.cache.setItem(key, arr)),
        map(
          arr =>
            ({ state: ApiState.LOADED, data: arr }) as ApiResponse<DateModel[]>
        ),
        startWith({ state: ApiState.LOADING } as ApiResponse<DateModel[]>),
        catchError(e => of(err<DateModel[]>(e)))
      );
  };
}
