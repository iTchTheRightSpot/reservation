import { ApiResponse, ApiState } from '@root/app.model';

export const tabledata = <T>(m: ApiResponse<T>) => {
  if (m.state === ApiState.LOADING)
    return Array.from({ length: 10 }).map(() => ({}));
  else if (m.state === ApiState.ERROR)
    return Array.from({ length: 1 }).map(() => ({}));
  return m.data!!;
};
