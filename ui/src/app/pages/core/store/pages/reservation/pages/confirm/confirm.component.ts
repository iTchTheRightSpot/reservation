import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ReservationService } from '@store/pages/reservation/reservation.service';
import { ReservationModel } from '@store/pages/reservation/reservation.model';
import { SummaryHolderComponent } from '@store/pages/reservation/shared/summary-holder.component';
import { Button } from 'primeng/button';
import {
  FormBuilder,
  FormControl,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Message } from 'primeng/message';
import { ApiResponse, ApiState } from '@root/app.model';
import { FloatLabel } from 'primeng/floatlabel';
import { Textarea } from 'primeng/textarea';
import { DatesService } from '@store/pages/reservation/pages/dates/dates.service';
import { ConfirmService } from '@store/pages/reservation/pages/confirm/confirm.service';
import { Subject, switchMap, tap } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { TIMEZONE } from '@root/app.util';
import { Router } from '@angular/router';
import { CORE_ROUTE } from '@root/app.routes';
import { STORE_ROUTE } from '@pages/core/core.routes';
import { SERVICE_TYPES_ROUTE } from '@store/pages/reservation/reservation.routes';
import { RESERVATION_ROUTE } from '@store/store.routes';
import { Divider } from 'primeng/divider';
import { ConfirmModel, ServiceTypeModel } from '@shared/model/shared.model';

@Component({
  selector: 'app-confirm',
  imports: [
    SummaryHolderComponent,
    Button,
    IconField,
    InputIcon,
    InputText,
    Message,
    ReactiveFormsModule,
    FloatLabel,
    Textarea,
    Divider
  ],
  templateUrl: './confirm.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConfirmComponent {
  private readonly confirmService = inject(ConfirmService);
  protected readonly reservationService = inject(ReservationService);
  private readonly datesService = inject(DatesService);
  private readonly router = inject(Router);

  protected readonly state = ApiState;

  protected readonly form = inject(FormBuilder).group({
    phone: new FormControl<string>('', [
      Validators.required,
      Validators.minLength(9),
      Validators.maxLength(20)
    ]),
    name: new FormControl<string>('', [
      Validators.required,
      Validators.minLength(1),
      Validators.maxLength(100)
    ]),
    email: new FormControl<string>('', [
      Validators.required,
      Validators.email,
      Validators.maxLength(320)
    ]),
    description: new FormControl<string>('', [Validators.maxLength(255)])
  });

  protected readonly check = (o: ReservationModel) =>
    !o.services || o.services.length < 1 || !o.staff
      ? undefined
      : { services: o.services, staff: o.staff };

  protected readonly subtotal = (s: ServiceTypeModel[] | undefined) =>
    !s ? 0 : s.map(s => s.price).reduce((acc, curr) => acc + curr, 0);

  protected readonly emit = new Subject<ConfirmModel>();

  protected readonly req = toSignal(
    this.emit.asObservable().pipe(
      switchMap(o =>
        this.confirmService.reserve(o).pipe(
          tap(obj => {
            if (obj.state === ApiState.LOADED) {
              Object.keys(this.form.controls).forEach(k =>
                this.form.get(k)?.setValue('')
              );
              this.datesService.clearCache();
              this.reservationService.clear();
              this.router.navigate([
                `${CORE_ROUTE}/${STORE_ROUTE}/${RESERVATION_ROUTE}/${SERVICE_TYPES_ROUTE}`
              ]);
            }
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly submit = () => {
    if (this.form.invalid) return;
    const state = this.reservationService.reservationState();
    this.emit.next({
      staff_id: state.staff?.staff_id || '',
      name: this.form.controls['name'].value || '',
      email: this.form.controls['email'].value || '',
      description: this.form.controls['description'].value || '',
      phone: `${this.form.controls['phone'].value}` || '',
      services: state.services?.map(s => s.name.trim()) || [],
      timezone: TIMEZONE,
      time: state.datetime || ''
    });
  };
}
