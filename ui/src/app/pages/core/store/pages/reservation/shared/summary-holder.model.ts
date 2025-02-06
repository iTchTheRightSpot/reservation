import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { StaffModel } from '@shared/model/shared.model';

export interface SummaryHolderModel {
  services: ServiceTypeModel[];
  staff: StaffModel;
}
