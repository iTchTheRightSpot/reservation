import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { environment } from '@env/environment';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { DateModel, DummyDates } from './dates.model';
import { ApiResponse, ApiState, err, TIMEZONE } from '@root/app.util';
import { Cache } from '@root/cache';

@Injectable({
  providedIn: 'root'
})
export class DatesService {
  private static readonly cache = new Cache<string, DateModel[]>();

  private readonly http = inject(HttpClient);
  private readonly service = inject(ReservationService);

  readonly dates = (date: Date) => {
    if (!environment.production)
      return of<ApiResponse<DateModel[]>>({
        state: ApiState.LOADED,
        data: DummyDates(date.getMonth(), date.getFullYear())
      });

    const obj = this.service.reservationState();
    if (!obj.services || obj.services.length < 1 || !obj.staff)
      return of<ApiResponse<DateModel[]>>({
        state: ApiState.ERROR,
        message: 'missing: service & staff'
      });

    const key = `${obj.staff.staff_id}_${obj.services.map(s => s.name).join('_')}_${date.getMonth()}_${date.getFullYear()}_${TIMEZONE}`;

    return DatesService.cache.getItem(key).pipe(
      switchMap(arr =>
        arr
          ? of<ApiResponse<DateModel[]>>({
              state: ApiState.LOADED,
              data: arr
            })
          : this.req(key, obj.staff!.staff_id, date)
      )
    );
  };

  private readonly req = (key: string, staffId: string, date: Date) => {
    let params = new HttpParams();
    params = params.append('staff_id', staffId);
    params = params.append('day', date.getDate());
    params = params.append('month', 1 + date.getMonth());
    params = params.append('year', date.getFullYear());
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
        catchError(e => of(err<DateModel[]>(e)))
      );
  };
}
