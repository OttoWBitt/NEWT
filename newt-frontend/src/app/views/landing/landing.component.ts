import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { LoginService } from 'src/app/services/views/login.service';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.scss']
})
export class LandingComponent implements OnInit {

  loginForm: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
    private loginService: LoginService,
  ) { }

  ngOnInit(): void {
    this.initForm();
  }

  initForm() {
    this.loginForm = this.formBuilder.group({
      username: [null , Validators.required],
      password: [null , Validators.required],
    });
  }

  signIn(){
    let formData = this.loginForm.getRawValue()
    this.loginService.login(formData.username, formData.password).subscribe(
      (response)=>{
        this.router.navigate([Views.artifacts.url])
      },
      (error)=>{
        console.log(error)
      }
    )
  }

  passwordRecovery(){
    console.log('User has forgotten his password')
  }
}
