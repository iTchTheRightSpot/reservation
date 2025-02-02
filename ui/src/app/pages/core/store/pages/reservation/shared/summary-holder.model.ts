import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { StaffModel } from '@store/pages/reservation/pages/staff/staff.model';

export interface SummaryHolderModel {
  services: ServiceTypeModel[];
  staff: StaffModel;
}
