import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Views } from 'src/app/util/views.enum';

@Component({
  selector: 'newt-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {

  user: string;

  constructor(
    private router: Router
  ) { }

  ngOnInit(): void {
    this.user = this.getUser()
  }

  getUser(){
    const sessionData = JSON.parse(localStorage.getItem('currentUser'));
    return sessionData.user.name
  }

  logout(){
    localStorage.clear()
    this.router.navigate([Views.login.url])
  }

}
