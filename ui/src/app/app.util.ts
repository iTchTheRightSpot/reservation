import { HttpErrorResponse } from '@angular/common/http';

export enum ApiState {
  LOADING = 'LOADING',
  LOADED = 'LOADED',
  ERROR = 'ERROR'
}

export interface ApiResponse<T> {
  data?: T;
  state: ApiState;
  message?: string;
}

export const err = <T>(e: HttpErrorResponse): ApiResponse<T> => ({
  state: ApiState.ERROR,
  message: (e.error ? e.error.message : e.message) as string
});
