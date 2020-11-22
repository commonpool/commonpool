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

  ngOnInit(): void {
  }

}
