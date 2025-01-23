import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CrmNavBarComponent } from '@crm/shared/crm-nav-bar.component';
import { AuthService } from '@shared/data-access/auth.service';
import { ProgressBar } from 'primeng/progressbar';
import { LoadingService } from '@shared/data-access/loading.service';

@Component({
  selector: 'app-staff',
  imports: [RouterOutlet, CrmNavBarComponent, ProgressBar],
  template: `
    <div class="w-full xl:max-w-7xl m-auto flex flex-col">
      <div class="sticky top-0 z-10">
        @if (loading.state()) {
          <p-progress-bar mode="indeterminate" [style]="{ height: '6px' }" />
        }
        <app-crm-nav-bar [imageKey]="auth.activeUser()?.image_key" />
      </div>
      <div class="py-3 px-2 xl:px-0 flex-1">
        <router-outlet />
      </div>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CrmComponent {
  protected readonly auth = inject(AuthService);
  protected readonly loading = inject(LoadingService);
}
