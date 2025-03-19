import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { Menubar } from 'primeng/menubar';
import { MenuItem } from 'primeng/api';
import { LightDarkModeComponent } from '@shared/ui/light-dark-mode.component';
import { StoreRoutes } from '@store/store.routes';
import { RootRoutes } from '@root/app.routes';
import { Avatar } from 'primeng/avatar';
import { CoreRoutes } from '@pages/core/core.routes';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-navigation',
  imports: [Menubar, LightDarkModeComponent, Avatar, RouterLink],
  template: `
    <div class="w-full">
      <p-menubar [model]="items">
        <ng-template #end>
          <div class="flex gap-2 items-center justify-center">
            @if (imageKey(); as k) {
              <a [routerLink]="crm">
                <p-avatar [image]="k" class="mr-2" shape="circle" />
              </a>
            } @else if (imageKey() === null) {
              <a [routerLink]="crm">
                <p-avatar icon="pi pi-user" class="mr-2" shape="circle" />
              </a>
            }
            <app-light-dark-mode />
          </div>
        </ng-template>
      </p-menubar>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class NavigationComponent {
  imageKey = input.required<string | null | undefined>();

  protected readonly crm = `${RootRoutes.CRM}`;
  protected readonly items: MenuItem[] = [
    {
      label: 'Home',
      icon: 'pi pi-home',
      routerLink: `/${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.HOME}`
    },
    {
      label: 'Book',
      icon: 'pi pi-calendar',
      routerLink: `/${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.RESERVATION}`
    },
    {
      label: 'About',
      icon: 'pi pi-search',
      routerLink: `/${RootRoutes.CORE}/${CoreRoutes.STORE}/${StoreRoutes.ABOUT}`
    },
    {
      label: 'Login',
      icon: 'pi pi-user',
      routerLink: `/${RootRoutes.CORE}/${CoreRoutes.LOGIN}`
    }
  ];
}
