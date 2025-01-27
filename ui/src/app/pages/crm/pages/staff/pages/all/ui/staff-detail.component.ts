import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { CRMStaffModel } from '@crm/pages/staff/pages/all/crm-staff.model';
import { Tab, TabList, TabPanel, TabPanels, Tabs } from 'primeng/tabs';
import { ViewStaffComponent } from './shared/view/view-staff.component';
import { Tag } from 'primeng/tag';
import { RoleComponent } from './shared/role/role.component';
import { LinkServiceComponent } from './shared/link-service/link-service.component';

@Component({
  selector: 'app-staff-detail-holder',
  imports: [
    Tabs,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
    ViewStaffComponent,
    Tag,
    RoleComponent,
    LinkServiceComponent
  ],
  template: `
    <p-tabs value="0">
      <p-tablist>
        <p-tab value="0">
          Account
          @if (staff().locked) {
            <p-tag [value]="'Locked'" [severity]="'danger'" />
          }
        </p-tab>
        <p-tab value="1">Roles & Permissions</p-tab>
        <p-tab value="2">Sevices</p-tab>
      </p-tablist>
      <p-tabpanels>
        <p-tabpanel value="0">
          <app-view-staff [staff]="staff()" />
        </p-tabpanel>
        <p-tabpanel value="1">
          <app-role [staff]="staff()" />
        </p-tabpanel>
        <p-tabpanel value="2">
          <app-link-service [staff]="staff()" />
        </p-tabpanel>
      </p-tabpanels>
    </p-tabs>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class StaffDetailComponent {
  staff = input.required<CRMStaffModel>();
}
