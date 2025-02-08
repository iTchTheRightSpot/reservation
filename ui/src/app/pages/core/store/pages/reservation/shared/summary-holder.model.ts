import { ServiceTypeModel, StaffModel } from '@shared/model/shared.model';

export interface SummaryHolderModel {
  services: ServiceTypeModel[];
  staff: StaffModel;
}
