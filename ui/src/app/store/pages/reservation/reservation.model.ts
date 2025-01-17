import { ServiceTypeModel } from './service-type/service-type.model';
import { StaffModel } from './staff/staff.model';

export interface ReservationModel {
  services: ServiceTypeModel[] | undefined;
  staff: StaffModel | undefined;
}
