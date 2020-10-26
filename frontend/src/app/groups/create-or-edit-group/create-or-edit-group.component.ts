import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {CreateGroupRequest} from '../../api/models';
import {Router} from '@angular/router';

@Component({
  selector: 'app-create-or-edit-group',
  templateUrl: './create-or-edit-group.component.html',
  styleUrls: ['./create-or-edit-group.component.css']
})
export class CreateOrEditGroupComponent implements OnInit {

  name = '';
  description = '';
  pending = false;
  error = undefined;

  constructor(private backend: BackendService, private router: Router) {

  }

  ngOnInit(): void {
  }

  submit() {
    this.pending = true;
    this.backend.createGroup(new CreateGroupRequest(this.name, this.description)).subscribe(res => {
      this.router.navigateByUrl('/');
      this.pending = false;
    }, err => {
      this.error = err;
      this.pending = false;
    });
  }

}
