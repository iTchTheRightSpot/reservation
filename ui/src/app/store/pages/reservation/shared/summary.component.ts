import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { Card } from 'primeng/card';
import { Avatar } from 'primeng/avatar';
import { ServiceTypeModel } from '@store/pages/reservation/pages/service-type/service-type.model';
import { ToggleButton } from 'primeng/togglebutton';
import { FormsModule } from '@angular/forms';
import { FORMAT_SECONDS } from '@root/app.util';
import { NgClass } from '@angular/common';

@Component({
  selector: 'app-summary',
  imports: [Card, Avatar, ToggleButton, FormsModule, NgClass],
  template: `
    <div>
      @if (displayHeader()) {
        <div class="mb-3">
          <h1 class="font-semibold text-base lg:text-lg">
            Appointment summary
          </h1>
        </div>
      }

      <div class="w-full">
        <p-card>
          <div class="m-0 flex gap-3">
            @if (staffImage(); as img) {
              <p-avatar
                [image]="img"
                class="mr-2"
                size="xlarge"
                shape="circle"
              />
            } @else {
              <p-avatar
                icon="pi pi-user"
                class="mr-2"
                size="large"
                shape="circle"
              />
            }

            <div class="flex-1">
              <h4 class="capitalize font-semibold">{{ staffName() }}</h4>
              <div class="mt2 flex justify-between">
                <p>{{ bio() }}</p>
                <p-toggle-button
                  [(ngModel)]="toggle"
                  onIcon="pi pi-chevron-up"
                  offIcon="pi pi-chevron-down"
                  onLabel="Up"
                  offLabel="Down"
                  size="small"
                  styleClass="text-center"
                />
              </div>
            </div>
          </div>
        </p-card>

        <ul
          class="max-h-96 m-0 p-0 pb-2 list-none overflow-y-auto"
          [ngClass]="{ hidden: toggle }"
        >
          @for (s of services(); track s.name + s.price + s.duration) {
            <li class="mt-4" tabindex="0">
              <p-card [header]="s.name">
                <p>
                  <span>
                    ₦{{ s.price }} <strong>.</strong>
                    {{ formatSeconds(s.duration) }}
                  </span>
                </p>
              </p-card>
            </li>
          }
        </ul>
      </div>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SummaryComponent {
  services = input.required<ServiceTypeModel[]>();
  staffImage = input.required<string | null>();
  staffName = input.required<string>();
  bio = input.required<string>();
  displayHeader = input<boolean>(true);

  protected toggle = false;

  protected readonly formatSeconds = (seconds: number) =>
    FORMAT_SECONDS(seconds);
}
