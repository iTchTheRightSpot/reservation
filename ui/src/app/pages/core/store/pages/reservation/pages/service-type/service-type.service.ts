import { inject, Injectable } from '@angular/core';
import { environment } from '@env/environment';
import { HttpClient } from '@angular/common/http';
import { DummyServiceTypes, ServiceTypeModel } from './service-type.model';
import { BehaviorSubject, catchError, map, of, switchMap, tap } from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';

@Injectable({
  providedIn: 'root'
})
export class ServiceTypeService {
  private readonly http = inject(HttpClient);

  private readonly cache = new BehaviorSubject<ServiceTypeModel[] | undefined>(
    undefined
  );

  readonly services = () =>
    !environment.production
      ? of<ApiResponse<ServiceTypeModel[]>>({
          data: DummyServiceTypes,
          state: ApiState.LOADED
        })
      : this.cache.asObservable().pipe(
          switchMap(arr =>
            arr
              ? of<ApiResponse<ServiceTypeModel[]>>({
                  data: arr,
                  state: ApiState.LOADED
                })
              : this.req$()
          )
        );

  private readonly req$ = () =>
    this.http
      .get<
        ServiceTypeModel[]
      >(`${environment.domain}service`, { withCredentials: true })
      .pipe(
        tap(arr => this.cache.next(arr)),
        map(
          arr =>
            ({ state: ApiState.LOADED, data: arr }) as ApiResponse<
              ServiceTypeModel[]
            >
        ),
        catchError(e => of(err<ServiceTypeModel[]>(e)))
      );
}
