from django.db import models

# Create your models here.
class Users(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.TextField('name', max_length=50)
    username = models.TextField('username',max_length=100)
    password = models.TextField('password',max_length=50)
    status = models.IntegerField('status',default=1)
    userClass = models.TextField('class')
    count = models.IntegerField('count',default=0)
    lastLogin = models.DateField('lastLogin')
    maxBet = models.IntegerField('maxBet',default=0)
    def __str__(self):
        return f'Id : {self.id}, Name: {self.name}, Count : {self.count}, Password: {self.password}, class {self.userClass}'
    def authenicate(request, username, password):
        try:
            return Users.objects.get(name = username, password = password)
        except Users.DoesNotExist:
            return None
    