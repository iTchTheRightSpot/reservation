import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { InputIcon } from 'primeng/inputicon';
import { IconField } from 'primeng/iconfield';
import { InputText } from 'primeng/inputtext';
import {
  FormBuilder,
  FormControl,
  FormsModule,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { Button } from 'primeng/button';
import { Message } from 'primeng/message';
import { Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { Subject, switchMap, tap } from 'rxjs';
import { ApiResponse, ApiState } from '@root/app.model';
import { STAFF_ROUTE } from '@root/app.routes';
import { ProgressSpinner } from 'primeng/progressspinner';
import { AuthService, LoginModel } from '@shared/data-access/auth.service';

@Component({
  selector: 'app-login',
  imports: [
    InputIcon,
    IconField,
    InputText,
    FormsModule,
    Button,
    Message,
    ReactiveFormsModule,
    ProgressSpinner
  ],
  templateUrl: 'login.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LoginComponent {
  // svg background link https://www.svgbackgrounds.com/set/free-svg-backgrounds-and-patterns/
  private readonly service = inject(AuthService);
  private readonly router = inject(Router);
  protected readonly form = inject(FormBuilder).group({
    email: new FormControl('', [Validators.required, Validators.email]),
    password: new FormControl('', [Validators.required])
  });

  protected readonly state = ApiState;

  private readonly subject = new Subject<LoginModel>();
  protected readonly login = toSignal(
    this.subject.asObservable().pipe(
      switchMap(obj => this.service.login(obj)),
      tap(obj => {
        if (obj.state === ApiState.LOADED)
          this.router.navigate([`${STAFF_ROUTE}`]);
      })
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<any> }
  );

  protected readonly submit = () => {
    if (this.form.invalid) return;
    this.subject.next(this.form.value as LoginModel);
  };
}
