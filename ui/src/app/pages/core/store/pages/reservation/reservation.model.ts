import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { StaffModel } from '@shared/model/shared.model';

export interface ReservationModel {
  services: ServiceTypeModel[] | undefined;
  staff: StaffModel | undefined;
  datetime: string | undefined;
}
