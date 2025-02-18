import {
  ChangeDetectionStrategy,
  Component,
  input,
  output
} from '@angular/core';
import { Avatar } from 'primeng/avatar';
import { LightDarkModeComponent } from '@shared/ui/light-dark-mode.component';
import { MenuItem } from 'primeng/api';
import { Menubar } from 'primeng/menubar';
import {
  SCHEDULE_ROUTE,
  SETTINGS_ROUTE
} from '@crm/pages/account/account.routes';
import {
  ACCOUNT_ROUTE,
  BOOKINGS_ROUTE,
  CRM_SERVICE_TYPES_ROUTE,
  CRM_STAFFS_ROUTE,
  DASHBOARD_ROUTE
} from '@crm/crm.routes';
import { CRM_ROUTE } from '@root/app.routes';
import {
  ALL_STAFFS_ROUTE,
  REGISTER_ROUTE
} from '@crm/pages/staff/staff-core.routes';
import { ApiState } from '@root/app.model';

@Component({
  selector: 'app-crm-nav-bar',
  imports: [Avatar, LightDarkModeComponent, Menubar],
  template: `
    <p-menubar [model]="items">
      <ng-template #end>
        <div class="flex gap-2 items-center justify-center">
          @if (imageKey(); as k) {
            <p-avatar [image]="k" class="mr-2" shape="circle" />
          } @else if (imageKey() === null) {
            <p-avatar icon="pi pi-user" class="mr-2" shape="circle" />
          }
          <app-light-dark-mode />
        </div>
      </ng-template>
    </p-menubar>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CrmNavBarComponent {
  imageKey = input.required<string | null | undefined>();
  logoutState = input.required<ApiState>();

  readonly logout = output<void>();

  protected readonly logo = './logo.jpeg';
  protected readonly items: MenuItem[] = [
    {
      label: 'Dashboard',
      icon: 'pi pi-gauge',
      routerLink: `/${CRM_ROUTE}/${DASHBOARD_ROUTE}`
    },
    {
      label: 'Bookings',
      icon: 'pi pi-calendar-clock',
      routerLink: `/${CRM_ROUTE}/${BOOKINGS_ROUTE}`
    },
    {
      label: 'Services',
      icon: 'pi pi-shop',
      routerLink: `/${CRM_ROUTE}/${CRM_SERVICE_TYPES_ROUTE}`
    },
    {
      label: 'Staffs',
      icon: 'pi pi-users',
      items: [
        {
          label: 'All',
          icon: 'pi pi-list',
          routerLink: `/${CRM_ROUTE}/${CRM_STAFFS_ROUTE}/${ALL_STAFFS_ROUTE}`
        },
        {
          label: 'Register',
          icon: 'pi pi-user-plus',
          routerLink: `/${CRM_ROUTE}/${CRM_STAFFS_ROUTE}/${REGISTER_ROUTE}`
        }
      ]
    },
    {
      label: 'Account',
      icon: 'pi pi-home',
      items: [
        {
          label: 'Schedule',
          icon: 'pi pi-calendar-clock',
          routerLink: `/${CRM_ROUTE}/${ACCOUNT_ROUTE}/${SCHEDULE_ROUTE}`
        },
        {
          label: 'Settings',
          icon: 'pi pi-spin pi-cog',
          routerLink: `/${CRM_ROUTE}/${ACCOUNT_ROUTE}/${SETTINGS_ROUTE}`
        },
        {
          label: 'Logout',
          icon: 'pi pi-power-off',
          command: () => this.logout.emit()
        }
      ]
    }
  ];
}
