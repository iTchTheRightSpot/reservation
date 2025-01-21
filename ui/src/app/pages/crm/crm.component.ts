import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CrmNavBarComponent } from '@crm/shared/crm-nav-bar.component';
import { AuthService } from '@shared/data-access/auth.service';

@Component({
  selector: 'app-staff',
  imports: [RouterOutlet, CrmNavBarComponent],
  template: `
    <div class="w-full xl:w-cx-50 m-auto flex flex-col">
      <div class="sticky top-0 z-10">
        <app-crm-nav-bar [imageKey]="service.activeUser()?.image_key" />
      </div>
      <div class="w-full py-3 px-2 xl:px-0 flex-1">
        <router-outlet />
      </div>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CrmComponent {
  protected readonly service = inject(AuthService);
}
