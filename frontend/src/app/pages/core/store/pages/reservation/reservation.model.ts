import { ServiceTypeModel, StaffModel } from '@shared/model/shared.model';

export interface ReservationModel {
  services: ServiceTypeModel[] | undefined;
  staff: StaffModel | undefined;
  datetime: string | undefined;
}
