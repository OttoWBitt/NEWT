import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/views/auth.service';
import { MessagesEnum } from 'src/app/util/messages.enum';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-password-recovery',
  templateUrl: './password-recovery.component.html',
  styleUrls: ['./password-recovery.component.scss']
})
export class PasswordRecoveryComponent implements OnInit {

  token: string
  passwordForm: FormGroup

  constructor(
    private formBuilder: FormBuilder,
    private authService: AuthService,
    private snackbar: MatSnackBar,
    private route: ActivatedRoute,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.token = this.route.snapshot.paramMap.get('token')
    this.initForm()
  }

  initForm(){
    this.passwordForm = this.formBuilder.group({
      password: [null, [Validators.required]],
      confirmPassword: [null, [Validators.required]],
    })
  }

  resetPassword(){
    if (this.passwordForm.valid && this.checkPasswords()){
      this.authService.resetPassword(this.token, this.passwordForm.get('password').value).subscribe(response => {
        if (response){
          this.snackbar.open(MessagesEnum.SuccessPasswordChange)
          this.router.navigate([Views.artifacts.url])
        }
      })
    }
  }

  checkPasswords() {
    const password = this.passwordForm.get('password').value;
    const confirmPassword = this.passwordForm.get('confirmPassword').value;
    this.passwordForm.controls['password'].setErrors(null)
    this.passwordForm.controls['confirmPassword'].setErrors(null)
    if (password == confirmPassword){
      return true
    } else {
      this.snackbar.open(MessagesEnum.passwordNotMatch, 'Fechar')
      this.passwordForm.controls['confirmPassword'].setErrors({'incorrect': true})
      return false
    }
  }
}
