import { HttpEvent, HttpHandlerFn, HttpRequest } from '@angular/common/http';
import { catchError, finalize, Observable, throwError } from 'rxjs';
import { inject } from '@angular/core';
import { LoadingService } from './loading.service';
import { Router } from '@angular/router';
import { CORE_ROUTE } from '@root/app.routes';
import { LOGIN_ROUTE } from '@pages/core/core.routes';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';

// responsible for loading ui and formatting error
export function loadingInterceptor(
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> {
  const service = inject(LoadingService).state;
  const router = inject(Router);
  const toast = inject(ToastService);
  service.set(true);
  return next(req).pipe(
    catchError(err => {
      const mess = err.error ? err.error.message : err.message;
      if (err.status === 401) {
        router.navigate([`${CORE_ROUTE}/${LOGIN_ROUTE}`]);
        toast.message({ message: mess, state: ToastEnum.ERROR });
      }
      return throwError(() => ({
        message: mess,
        status: err.status
      }));
    }),
    finalize(() => service.set(false))
  );
}
