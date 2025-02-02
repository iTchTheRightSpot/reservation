import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { Button } from 'primeng/button';
import { PrimeTemplate } from 'primeng/api';
import { SummaryComponent } from './summary.component';
import { SummaryHolderModel } from './summary-holder.model';
import { Drawer } from 'primeng/drawer';

@Component({
  selector: 'app-summary-holder',
  imports: [Button, PrimeTemplate, SummaryComponent, Drawer],
  template: `
    <div class="w-full hidden md:block">
      <app-summary
        [staffImage]="obj().staff.image_key"
        [staffName]="obj().staff.name"
        [bio]="obj().staff.bio"
        [services]="obj().services"
        [datetime]="datetime()"
      />
    </div>

    <div class="w-full block md:hidden fixed left-0 right-0 bottom-0">
      <p-drawer
        styleClass="!h-[30rem]"
        [(visible)]="sidebarVisible4"
        position="bottom"
        [dismissible]="false"
      >
        <ng-template pTemplate="header">
          <span class="font-semibold text-xl">Appointment Summary</span>
        </ng-template>

        <ng-template pTemplate="content">
          <app-summary
            [staffImage]="obj().staff.image_key"
            [staffName]="obj().staff.name"
            [bio]="obj().staff.bio"
            [services]="obj().services"
            [datetime]="datetime()"
            [displayHeader]="false"
          />
        </ng-template>
      </p-drawer>

      <div class="w-full p-2 text-right">
        <p-button
          type="button"
          icon="pi pi-arrow-up"
          (onClick)="sidebarVisible4 = true"
        />
      </div>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SummaryHolderComponent {
  protected sidebarVisible4 = false;

  obj = input.required<SummaryHolderModel>();
  datetime = input.required<string | undefined>();
}
