import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { ToastEnum, ToastService } from '@shared/data-access/toast.service';
import { tap } from 'rxjs';
import { MessageService } from 'primeng/api';
import { Toast } from 'primeng/toast';
import { AsyncPipe } from '@angular/common';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, Toast, AsyncPipe],
  providers: [MessageService],
  template: `
    <p-toast />
    @if (toast$ | async) {}
    <router-outlet />
  `,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AppComponent {
  private readonly service = inject(MessageService);

  protected readonly toast$ = inject(ToastService).toast$.pipe(
    tap(obj => {
      if (obj.state === ToastEnum.ERROR)
        this.service.add({
          severity: 'error',
          summary: 'Error',
          detail: obj.message
        });
      else if (obj.state === ToastEnum.SUCCESS)
        this.service.add({
          severity: 'success',
          summary: 'Success',
          detail: obj.message
        });
    })
  );
}
