import * as moment from 'moment-timezone';
import { HttpErrorResponse } from '@angular/common/http';
import { ApiResponse, ApiState } from '@root/app.model';
import { DateModel } from '@shared/model/shared.model';

export const err = <T>(e: HttpErrorResponse): ApiResponse<T> => ({
  state: ApiState.ERROR,
  message: e.message
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
