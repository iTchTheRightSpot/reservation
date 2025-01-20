import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import {
  NavigationEnd,
  Router,
  RouterLink,
  RouterOutlet
} from '@angular/router';
import { Breadcrumb } from 'primeng/breadcrumb';
import { MenuItem } from 'primeng/api';
import { RESERVATION_ROUTE } from '@store/store.routes';
import {
  CONFIRM_ROUTE,
  DATES_ROUTE,
  SERVICE_TYPE_ROUTE,
  STORE_STAFF_ROUTE
} from '@store/pages/reservation/reservation.routes';
import { NgClass } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';

@Component({
  selector: 'app-reservation',
  imports: [RouterOutlet, Breadcrumb, RouterLink, NgClass],
  template: ` <div class="w-full xl:w-cx-50 m-auto">
    <div class="mt-3">
      <p-breadcrumb [model]="items" [home]="home">
        <ng-template #item let-item>
          <a [routerLink]="item.route" class="p-breadcrumb-item-link">
            <span
              class="text-primary font-semibold"
              [ngClass]="{ 'text-primary': item.route === url() }"
              >{{ item.label }}</span
            >
          </a>
        </ng-template>
      </p-breadcrumb>
    </div>

    <div class="pb-3 xl:pb-0 px-3 xl:px-0">
      <router-outlet />
    </div>
  </div>`,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ReservationComponent {
  private readonly router = inject(Router);

  protected readonly url = toSignal(
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(event => event.url),
      startWith(this.router.url)
    ),
    { initialValue: `/${RESERVATION_ROUTE}/${SERVICE_TYPE_ROUTE}` }
  );

  protected readonly items: MenuItem[] = [
    {
      label: 'services',
      route: `/${RESERVATION_ROUTE}/${SERVICE_TYPE_ROUTE}`
    },
    {
      label: 'staffs',
      route: `/${RESERVATION_ROUTE}/${STORE_STAFF_ROUTE}`
    },
    {
      label: 'dates',
      route: `/${RESERVATION_ROUTE}/${DATES_ROUTE}`
    },
    {
      label: 'confirm',
      route: `/${RESERVATION_ROUTE}/${CONFIRM_ROUTE}`
    }
  ];

  protected readonly home: MenuItem = {
    label: 'reservation',
    icon: 'pi pi-shopping-bag',
    routerLink: `/${RESERVATION_ROUTE}/${SERVICE_TYPE_ROUTE}`
  };
}
