import { RolePermissionEntity } from '@crm/pages/staff/pages/all/crm-staff.model';

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

export interface ActiveUser {
  user_id: string;
  firstname: string;
  image_key: string | null;
  access_controls: RolePermissionEntity[];
}
