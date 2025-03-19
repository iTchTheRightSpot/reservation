import { inject, Injectable } from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Cache } from '@shared/data-access/cache';
import { environment } from '@env/environment';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { DummyStaffModels, StaffModel } from '@shared/model/shared.model';

@Injectable({
  providedIn: 'root'
})
export class StaffService {
  private static readonly cache = new Cache<string, StaffModel[]>();

  private readonly http = inject(HttpClient);
  private readonly service = inject(ReservationService);

  readonly staffs = () => {
    if (!environment.production)
      return of<ApiResponse<StaffModel[]>>({
        state: ApiState.LOADED,
        data: DummyStaffModels
      });

    const services = this.service.reservationState().services;

    if (!services)
      return of<ApiResponse<StaffModel[]>>({
        state: ApiState.ERROR,
        message: 'please select 1 or more services to pre-book'
      });

    const key = services.map(s => s.name).join('_');

    return StaffService.cache.getItem(key).pipe(
      switchMap(value => {
        if (value)
          return of<ApiResponse<StaffModel[]>>({
            state: ApiState.LOADED,
            data: value
          });

        let params = new HttpParams();
        services.forEach(
          service => (params = params.append('name', service.name.trim()))
        );

        return this.req(params, key);
      })
    );
  };

  private readonly req = (params: HttpParams, key: string) =>
    this.http
      .get<
        StaffModel[]
      >(`${environment.domain}service/staffs`, { withCredentials: true, params: params })
      .pipe(
        tap(staffs => StaffService.cache.setItem(key, staffs)),
        map(
          arr =>
            ({ state: ApiState.LOADED, data: arr }) as ApiResponse<StaffModel[]>
        ),
        catchError(e => of(err<StaffModel[]>(e)))
      );
}
