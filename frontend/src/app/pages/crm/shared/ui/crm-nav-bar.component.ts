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
import { AccountRoutes } from '@crm/pages/account/account.routes';
import { CRMRoutes } from '@crm/crm.routes';
import { RootRoutes } from '@root/app.routes';
import { StaffRoutes } from '@crm/pages/staff/staff-core.routes';
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

  protected readonly items: MenuItem[] = [
    {
      label: 'Dashboard',
      icon: 'pi pi-gauge',
      routerLink: `/${RootRoutes.CRM}/${CRMRoutes.DASHBOARD}`
    },
    {
      label: 'Bookings',
      icon: 'pi pi-calendar-clock',
      routerLink: `/${RootRoutes.CRM}/${CRMRoutes.BOOKINGS}`
    },
    {
      label: 'Services',
      icon: 'pi pi-shop',
      routerLink: `/${RootRoutes.CRM}/${CRMRoutes.SERVICES}`
    },
    {
      label: 'Staffs',
      icon: 'pi pi-users',
      items: [
        {
          label: 'All',
          icon: 'pi pi-list',
          routerLink: `/${RootRoutes.CRM}/${CRMRoutes.STAFFS}/${StaffRoutes.STAFFS}`
        },
        {
          label: 'Register',
          icon: 'pi pi-user-plus',
          routerLink: `/${RootRoutes.CRM}/${CRMRoutes.STAFFS}/${StaffRoutes.REGISTER}`
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
          routerLink: `/${RootRoutes.CRM}/${CRMRoutes.ACCOUNT}/${AccountRoutes.SCHEDULE}`
        },
        {
          label: 'Settings',
          icon: 'pi pi-spin pi-cog',
          routerLink: `/${RootRoutes.CRM}/${CRMRoutes.ACCOUNT}/${AccountRoutes.SETTINGS}`
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
