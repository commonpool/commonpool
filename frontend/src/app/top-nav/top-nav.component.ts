import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {BackendService} from '../api/backend.service';

@Component({
  selector: 'app-top-nav',
  templateUrl: './top-nav.component.html',
  styleUrls: ['./top-nav.component.css']
})
export class TopNavComponent implements OnInit {

  constructor(public auth: AuthService, public backend: BackendService) {
  }

  bla = `:::user:::bfac680d-37c4-45a5-b4d0-ecd136958016:::`;

  ngOnInit(): void {
  }

}
