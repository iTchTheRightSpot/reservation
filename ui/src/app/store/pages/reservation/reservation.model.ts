import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { StaffModel } from '@store/pages/reservation/pages/staff/staff.model';

export interface ReservationModel {
  services: ServiceTypeModel[] | undefined;
  staff: StaffModel | undefined;
}
