import { HttpErrorResponse } from '@angular/common/http';
import * as moment from 'moment-timezone';

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

export const FORMAT_SECONDS = (seconds: number) =>
  `${hrImpl(Math.floor(seconds / 3600))} ${minImpl(Math.floor((seconds % 3600) / 60))} ${secImpl(seconds % 60)}`;

const hrImpl = (utcHours: number) => {
  if (utcHours < 1) {
    return '';
  } else if (utcHours === 1) {
    return '1 hr ';
  }
  return `${utcHours} hrs `;
};

const minImpl = (utcMins: number) => {
  if (utcMins < 1) {
    return '';
  } else if (utcMins === 1) {
    return '1 min ';
  }
  return `${utcMins} mins `;
};

const secImpl = (seconds: number) => {
  if (seconds < 1) {
    return '';
  } else if (seconds === 1) {
    return '1 sec';
  }
  return `${seconds} secs`;
};

export const TIMEZONE = moment.tz.guess();
