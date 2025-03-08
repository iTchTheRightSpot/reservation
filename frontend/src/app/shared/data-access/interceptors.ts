import { HttpEvent, HttpHandlerFn, HttpRequest } from '@angular/common/http';
import { catchError, finalize, Observable, throwError } from 'rxjs';
import { inject } from '@angular/core';
import { LoadingService } from './loading.service';
import { Router } from '@angular/router';
import { RootRoutes } from '@root/app.routes';
import { CoreRoutes } from '@pages/core/core.routes';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';

// responsible for loading ui and formatting error
export function interceptor(
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
      if (err.status === 401 && !req.url.endsWith('/api/v1/active')) {
        router.navigate([`${RootRoutes.CORE}/${CoreRoutes.LOGIN}`]);
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
