// Package conenv contains some tools to load settings from environment variable.
//
// conenv stands for "container" (and "configuration") with "environment variable".
// It aims to be feature-rich, so speed is not main concern.
//
// It is initially designed for dynamically scalable applications running in
// container:
//
//   1. Sometimes we have to pass dangerous data (like db password) to application.
//   2. Using configuration file is not best solution, as complex setup is required
//      when application need to be dynamic-scalable across multiple machines.
//   3. Loading struct value from environment variable is painful in Go.
//   4. There are some tools to help you load value with ease, you just have to
//      implement your own encrypting process, which is also a pain in the butt.
//
// conenv provides plugable support for data encryption, see DESExtension() for
// an example usecase.
//
// Here are tips about security when using conenv with your containerized program:
//
//   1. Encrypt sensitive info, password in particular.
//   2. Keep your keys in secret.
//   3. Sign sensitive info to prevent hackers from modifying it.
//   4. Embed neccessary keys in your program. Neither passing from external source,
//      nor embedding whole key.
//
// Take a look at package example code, which implements tips above at lowest level.
package conenv
