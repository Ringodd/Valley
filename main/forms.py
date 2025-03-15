from .models import Users
from django.forms import ModelForm, TextInput

class LoginForm(ModelForm):
    class Meta:
        model = Users
        fields = ['name', 'password']

        widgets = {
            'name': TextInput(attrs={
                'class':'form-control login',
                'placeholder':'Имя',
                'id':'form',
                'label':'username',
                'onchange':'checkInput()'
                }),
            'password':TextInput(attrs={
                'class':'form-control login',
                'placeholder':'Пароль',
                'type':'password',
                'id':'password',
                'onchange':'checkInput()'
            })
        }

