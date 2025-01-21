import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Button } from 'primeng/button';
import {
  FormBuilder,
  FormControl,
  FormsModule,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { IconField } from 'primeng/iconfield';
import { InputIcon } from 'primeng/inputicon';
import { InputText } from 'primeng/inputtext';
import { Message } from 'primeng/message';
import { ApiResponse, ApiState } from '@root/app.model';
import { AuthService, RegisterModel } from '@shared/data-access/auth.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { Subject, switchMap, tap } from 'rxjs';

@Component({
  selector: 'app-register',
  imports: [
    Button,
    FormsModule,
    IconField,
    InputIcon,
    InputText,
    Message,
    ReactiveFormsModule
  ],
  templateUrl: './register.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RegisterComponent {
  private readonly service = inject(AuthService);

  // regex link https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
  // min 8, max 15 characters, at least one uppercase letter, one lowercase letter, one number and one special character
  private readonly rx =
    '^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,15}$';
  protected readonly form = inject(FormBuilder).group({
    firstname: new FormControl('', [
      Validators.required,
      Validators.minLength(1),
      Validators.maxLength(50)
    ]),
    lastname: new FormControl('', [
      Validators.required,
      Validators.minLength(1),
      Validators.maxLength(50)
    ]),
    email: new FormControl('', [
      Validators.required,
      Validators.email,
      Validators.maxLength(320)
    ]),
    password: new FormControl('', [
      Validators.required,
      Validators.minLength(8),
      Validators.maxLength(15),
      Validators.pattern(this.rx)
    ]),
    confirm_password: new FormControl('', [
      Validators.required,
      Validators.minLength(8),
      Validators.maxLength(15),
      Validators.pattern(this.rx)
    ])
  });

  protected readonly state = ApiState;
  protected viewPassword = false;

  private readonly sub = new Subject<RegisterModel>();

  protected readonly register = toSignal(
    this.sub.asObservable().pipe(
      switchMap(o =>
        this.service.register(o).pipe(
          tap(obj => {
            if (obj.state === ApiState.LOADED)
              Object.keys(this.form.controls).forEach(key =>
                this.form.get(key)?.reset('')
              );
          })
        )
      )
    ),
    { initialValue: { state: ApiState.LOADED } as ApiResponse<RegisterModel> }
  );

  // validates if password match
  protected readonly match = () =>
    this.form.controls['password'].value ===
    this.form.controls['confirm_password'].value;

  protected readonly submit = () => {
    if (this.form.invalid) return;
    const { confirm_password, ...obj } = this.form.value;
    this.sub.next(obj as RegisterModel);
  };
}
