import { Component, OnInit } from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-home-page',
  templateUrl: './home-page.component.html',
  styleUrls: ['./home-page.component.css']
})
export class HomePageComponent implements OnInit {
  constructor(public backend: BackendService, public authService: AuthService) { }
  ngOnInit(): void {
  }
}

