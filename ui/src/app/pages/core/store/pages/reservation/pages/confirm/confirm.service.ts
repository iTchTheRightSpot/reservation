import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '@env/environment';
import { catchError, concat, concatMap, map, of, startWith, timer } from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { err } from '@root/app.util';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { ConfirmModel } from '@shared/data-access/shared.model';

@Injectable({
  providedIn: 'root'
})
export class ConfirmService {
  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly reserve = (c: ConfirmModel) =>
    environment.production
      ? this.http
          .post<any>(`${environment.domain}reservation`, c, {
            withCredentials: true
          })
          .pipe(
            map(() => {
              this.toast.message({
                state: ToastEnum.SUCCESS,
                message: 'appointment created!'
              });
              return { state: ApiState.LOADED } as ApiResponse<any>;
            }),
            startWith({ state: ApiState.LOADING } as ApiResponse<any>),
            catchError(e => of(err<any>(e)))
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
