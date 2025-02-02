import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { CrmNavBarComponent } from '@crm/shared/ui/crm-nav-bar.component';
import { AuthService } from '@shared/data-access/auth.service';
import { ProgressBar } from 'primeng/progressbar';
import { LoadingService } from '@shared/data-access/loading.service';
import { Subject, switchMap, tap } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ApiResponse, ApiState } from '@root/app.model';
import { CORE_ROUTE } from '@root/app.routes';
import { LOGIN_ROUTE } from '@pages/core/core.routes';

@Component({
  selector: 'app-staff',
  imports: [RouterOutlet, CrmNavBarComponent, ProgressBar],
  template: `
    <div class="w-full xl:max-w-7xl m-auto flex flex-col">
      <div class="sticky top-0 z-10">
        @if (loading.state()) {
          <p-progress-bar mode="indeterminate" [style]="{ height: '6px' }" />
        }
        <app-crm-nav-bar
          [imageKey]="auth.activeUser()?.image_key"
          [logoutState]="logout().state"
          (logout)="logoutEmitter.next($event)"
        />
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
  private readonly router = inject(Router);

  protected readonly logoutEmitter = new Subject<void>();
  protected readonly logout = toSignal(
    this.logoutEmitter.asObservable().pipe(
      switchMap(() =>
        this.auth.logout().pipe(
          tap(s => {
            if (s.state === ApiState.LOADED)
              this.router.navigate([`${CORE_ROUTE}/${LOGIN_ROUTE}`]);
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );
}
