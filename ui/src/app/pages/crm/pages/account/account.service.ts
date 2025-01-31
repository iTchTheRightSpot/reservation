import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { DeletePermission, DeleteRole, RoleAndPermissionPayload } from './account.model';
import { environment } from '@env/environment';
import { ApiResponse, ApiState } from '@root/app.model';
import { catchError, delay, map, merge, of, startWith } from 'rxjs';
import { err } from '@root/app.util';

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  private readonly http = inject(HttpClient);
  private readonly toast = inject(ToastService);

  readonly createRoleAndPermission = (o: RoleAndPermissionPayload) =>
    environment.production
      ? this.http
        .post<
          ApiResponse<any>
        >(`${environment.domain}account/role-permission`, o, { withCredentials: true })
        .pipe(
          map(() => <ApiResponse<any>>{ state: ApiState.LOADED }),
          startWith(<ApiResponse<any>>{ state: ApiState.LOADING }),
          catchError(e => {
            this.toast.message({
              message: e.message,
              state: ToastEnum.ERROR
            });
            return of(err<any>(e));
          })
        )
      : merge(
        of({ state: ApiState.LOADING }),
        of({ state: ApiState.LOADED }).pipe(delay(2000))
      );

  readonly deleteRole = (o: DeleteRole) =>
    environment.production
      ? this.http
        .delete<
          ApiResponse<any>
        >(`${environment.domain}account/role/${o.user_id}/${o.role}`, { withCredentials: true })
        .pipe(
          map(() => <ApiResponse<any>>{ state: ApiState.LOADED }),
          startWith(<ApiResponse<any>>{ state: ApiState.LOADING }),
          catchError(e => {
            this.toast.message({
              message: e.message,
              state: ToastEnum.ERROR
            });
            return of(err<any>(e));
          })
        )
      : merge(
        of({ state: ApiState.LOADING }),
        of({ state: ApiState.LOADED }).pipe(delay(2000))
      );

  readonly deletePermission = (o: DeletePermission) =>
    environment.production
      ? this.http
        .delete<
          ApiResponse<any>
        >(`${environment.domain}account/permission/${o.user_id}/${o.role}/${o.permission}`, { withCredentials: true })
        .pipe(
          map(() => <ApiResponse<any>>{ state: ApiState.LOADED }),
          startWith(<ApiResponse<any>>{ state: ApiState.LOADING }),
          catchError(e => {
            this.toast.message({
              message: e.message,
              state: ToastEnum.ERROR
            });
            return of(err<any>(e));
          })
        )
      : merge(
        of({ state: ApiState.LOADING }),
        of({ state: ApiState.LOADED }).pipe(delay(2000))
      );
}
