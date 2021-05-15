import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { User } from 'src/app/models/user.model';
import { AuthService } from 'src/app/services/views/auth.service';
import { generalExceptionTreatment } from 'src/app/util/error-handler';
import { MessagesEnum } from 'src/app/util/messages.enum';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.scss']
})
export class SignUpComponent implements OnInit {

  signUpForm: FormGroup;
  newUser: User = new User()

  constructor(
    private authService: AuthService,
    private snackbar: MatSnackBar,
    private router: Router,
    private formBuilder: FormBuilder
  ) { }

  ngOnInit(): void {
    this.initForm();
  }

  initForm() {
    this.signUpForm = this.formBuilder.group({
      name: [null , Validators.required],
      username: [null , Validators.required],
      password: [null , Validators.required],
      email: [null , Validators.required],
    })
  }

  signUp(){
    if (this.signUpForm.valid){
      this.setUserData()
      this.authService.signUp(this.newUser).subscribe(
        response => {
          if (response){
            this.snackbar.open(MessagesEnum.SuccessUserSignUp);
            this.router.navigate([Views.artifacts.url])
          }
        },
        error => {
          this.snackbar.open(generalExceptionTreatment(error), 'Fechar')
        }
      )
    }
  }

  setUserData(){
    let formData = this.signUpForm.getRawValue()
    this.newUser.name = formData.name
    this.newUser.username = formData.username
    this.newUser.password = formData.password
    this.newUser.email = formData.email
  }

  goToHome(){
    this.router.navigate([Views.login.url])
  }
}
