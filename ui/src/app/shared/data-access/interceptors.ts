import { HttpEvent, HttpHandlerFn, HttpRequest } from '@angular/common/http';
import { catchError, finalize, Observable, throwError } from 'rxjs';
import { inject } from '@angular/core';
import { LoadingService } from './loading.service';

// responsible for loading ui and formatting error
export function loadingInterceptor(
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<unknown>> {
  const service = inject(LoadingService).state;
  service.set(true);
  return next(req).pipe(
    catchError(err =>
      throwError(() => ({
        message: err.error ? err.error.message : err.message,
        status: err.status
      }))
    ),
    finalize(() => service.set(false))
  );
}
